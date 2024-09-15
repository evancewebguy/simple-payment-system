package database

import (
	"context"
	"fmt"
	"log"
	"mamlaka/config"
	"mamlaka/internal/app/payment"
	"mamlaka/internal/app/user"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetDB returns the GORM DB instance.
	GetDB() *gorm.DB
}

type service struct {
	DB *gorm.DB
}

var dbInstance *service

// New initializes a new database service and reuses the connection if already established.
func New(conf config.Config) Service {
	// Reuse connection
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		conf.Postgres.Host, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBName, conf.Postgres.Port, "disable", "Africa/Nairobi")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	dbInstance = &service{
		DB: db,
	}
	//Automatically Migrate
	err = db.AutoMigrate(
		user.User{},              // User model
		payment.Payment{},        // Payment model
		payment.PaymentDetails{}, // PaymentDetails model
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Database health check
	sqlDB, err := s.DB.DB()
	if err != nil {
		stats["db_status"] = "down"
		stats["db_error"] = fmt.Sprintf("failed to retrieve SQL DB from GORM: %v", err)
		log.Printf("failed to retrieve SQL DB from GORM: %v", err)
		return stats
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		stats["db_status"] = "down"
		stats["db_error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	// Database is up
	stats["db_status"] = "up"
	stats["db_message"] = "Database is healthy"

	// Database connection pool stats
	dbStats := sqlDB.Stats()
	stats["open_connections"] = fmt.Sprintf("%d", dbStats.OpenConnections)
	stats["in_use"] = fmt.Sprintf("%d", dbStats.InUse)
	stats["idle"] = fmt.Sprintf("%d", dbStats.Idle)
	stats["wait_count"] = fmt.Sprintf("%d", dbStats.WaitCount)
	stats["wait_duration"] = dbStats.WaitDuration.String()

	// System Health Checks (CPU, Memory, Disk)

	// CPU Info
	cpuPercent, _ := cpu.Percent(0, false)
	stats["cpu_usage_percent"] = fmt.Sprintf("%.2f%%", cpuPercent[0])

	// Memory Info
	vMem, _ := mem.VirtualMemory()
	stats["total_memory"] = strconv.FormatUint(vMem.Total/1024/1024/1024, 10) + " GB"
	stats["used_memory"] = strconv.FormatUint(vMem.Used/1024/1024/1024, 10) + " GB"
	stats["free_memory"] = strconv.FormatUint(vMem.Free/1024/1024/1024, 10) + " GB"
	stats["memory_usage_percent"] = fmt.Sprintf("%.2f%%", vMem.UsedPercent)

	// Disk Info
	diskUsage, _ := disk.Usage("/")
	stats["total_disk"] = strconv.FormatUint(diskUsage.Total/1024/1024/1024, 10) + " GB"
	stats["used_disk"] = strconv.FormatUint(diskUsage.Used/1024/1024/1024, 10) + " GB"
	stats["free_disk"] = strconv.FormatUint(diskUsage.Free/1024/1024/1024, 10) + " GB"
	stats["disk_usage_percent"] = fmt.Sprintf("%.2f%%", diskUsage.UsedPercent)

	// Host Info
	hostInfo, _ := host.Info()
	stats["host_uptime"] = fmt.Sprintf("%d minutes", hostInfo.Uptime/60)
	stats["host_os"] = hostInfo.OS
	stats["host_platform"] = hostInfo.Platform
	stats["host_kernel_version"] = hostInfo.KernelVersion

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	log.Println("Disconnecting from the database")
	return sqlDB.Close()
}

// GetDB returns the GORM DB instance.
func (s *service) GetDB() *gorm.DB {
	return s.DB
}
