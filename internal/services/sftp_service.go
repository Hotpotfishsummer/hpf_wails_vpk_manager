package services

import (
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type SftpService struct {
	mu         sync.Mutex
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	lastConfig models.SftpSettings
}

func NewSftpService() *SftpService {
	return &SftpService{}
}

func (s *SftpService) getClient(config models.SftpSettings) (*sftp.Client, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if config.Address == "" || config.Username == "" {
		return nil, fmt.Errorf("SFTP 地址或用户名未配置")
	}

	// If already connected with same config, verify connection
	if s.sftpClient != nil && s.lastConfig == config {
		if _, err := s.sftpClient.Getwd(); err == nil {
			return s.sftpClient, nil
		}
	}

	s.closeInternal()

	host := config.Address
	port := 22
	if strings.Contains(config.Address, ":") {
		parts := strings.Split(config.Address, ":")
		host = parts[0]
		p, err := strconv.Atoi(parts[1])
		if err != nil || p < 1 || p > 65535 {
			return nil, fmt.Errorf("无效的 SFTP 端口: %s", parts[1])
		}
		port = p
	}

	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	sshConn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("连接 SSH 服务器失败: %w", err)
	}

	sftpConn, err := sftp.NewClient(sshConn)
	if err != nil {
		_ = sshConn.Close()
		return nil, fmt.Errorf("初始化 SFTP 会话失败: %w", err)
	}

	s.sshClient = sshConn
	s.sftpClient = sftpConn
	s.lastConfig = config
	return s.sftpClient, nil
}

func (s *SftpService) closeInternal() {
	if s.sftpClient != nil {
		_ = s.sftpClient.Close()
		s.sftpClient = nil
	}
	if s.sshClient != nil {
		_ = s.sshClient.Close()
		s.sshClient = nil
	}
}

func (s *SftpService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeInternal()
}

func (s *SftpService) ListVpkFiles(config models.SftpSettings) ([]models.RemoteFile, error) {
	if config.AddonsPath == "" {
		return nil, fmt.Errorf("未配置远程 addons 路径")
	}

	client, err := s.getClient(config)
	if err != nil {
		return nil, err
	}

	entries, err := client.ReadDir(config.AddonsPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取远程目录 %s: %w", config.AddonsPath, err)
	}

	var files []models.RemoteFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		lower := strings.ToLower(name)
		if strings.HasSuffix(lower, ".vpk") || strings.HasSuffix(lower, ".vpk.disabled") {
			sizeMB := math.Round((float64(entry.Size())/(1024*1024))*100) / 100
			files = append(files, models.RemoteFile{
				Name:         name,
				Path:         path.Join(config.AddonsPath, name),
				Size:         sizeMB,
				ModifiedTime: float64(entry.ModTime().Unix()),
				IsDisabled:   strings.HasSuffix(lower, ".disabled"),
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return files, nil
}

func (s *SftpService) DeleteFile(config models.SftpSettings, remotePath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}
	return client.Remove(remotePath)
}

func (s *SftpService) DeleteFileWithImages(config models.SftpSettings, remotePath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}

	if err := client.Remove(remotePath); err != nil {
		return err
	}

	base := strings.TrimSuffix(remotePath, path.Ext(remotePath))
	for _, ext := range naming.AssociatedSyncExtensions {
		imgPath := base + ext
		if _, err := client.Stat(imgPath); err == nil {
			_ = client.Remove(imgPath)
		}
	}
	return nil
}

func (s *SftpService) RenameFile(config models.SftpSettings, oldPath, newPath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}
	return client.Rename(oldPath, newPath)
}

func (s *SftpService) RenameFileWithImages(config models.SftpSettings, oldPath, newPath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}

	if err := client.Rename(oldPath, newPath); err != nil {
		return err
	}

	oldBase := strings.TrimSuffix(oldPath, path.Ext(oldPath))
	newBase := strings.TrimSuffix(newPath, path.Ext(newPath))

	for _, ext := range naming.AssociatedSyncExtensions {
		oldImg := oldBase + ext
		newImg := newBase + ext
		if _, err := client.Stat(oldImg); err == nil {
			_ = client.Rename(oldImg, newImg)
		}
	}
	return nil
}

func (s *SftpService) FileExists(config models.SftpSettings, remotePath string) bool {
	client, err := s.getClient(config)
	if err != nil {
		return false
	}
	_, err = client.Stat(remotePath)
	return err == nil
}

func (s *SftpService) DownloadFile(config models.SftpSettings, remotePath, localPath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(localPath), 0755)

	srcFile, err := client.Open(remotePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (s *SftpService) UploadFile(config models.SftpSettings, localPath, remotePath string) error {
	client, err := s.getClient(config)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remotePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
