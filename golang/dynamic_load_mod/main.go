package main

import (
	"dynamic_load_mod/middleware"
	"fmt"
	"log"
)

func main() {

	// 读取配置文件，决定使用哪些中间件
	config := map[string]bool{
		"redis":      true,
		"kafka":      true,
		"clickhouse": true,
	}

	// 创建中间件管理器
	manager := middleware.NewManager()

	// 根据配置初始化中间件
	for mw, enabled := range config {
		if enabled {
			if err := manager.InitMiddleware(mw); err != nil {
				log.Printf("Failed to initialize %s: %v", mw, err)
			} else {
				log.Printf("%s initialized successfully", mw)
			}
		}
	}

	// 使用中间件
	if setter, ok := manager.GetWriter("redis"); ok {
		if err := setter.Write("key", "value"); err != nil {
			log.Printf("Failed to set key in Redis: %v", err)
		}
	} else {
		fmt.Println("Redis writer is not available")
	}

	if reader, ok := manager.GetReader("redis"); ok {
		if _, err := reader.Read("key"); err != nil {
			log.Printf("Failed to get key in Redis: %v", err)
		}
	} else {
		fmt.Println("Redis reader is not available")
	}

	if producer, ok := manager.GetWriter("kafka"); ok {
		if err := producer.Write("topic", "message"); err != nil {
			log.Printf("Failed to produce message to Kafka: %v", err)
		}
	} else {
		fmt.Println("Kafka producer is not available")
	}

	if selector, ok := manager.GetReader("clickhouse"); ok {
		if value, err := selector.Read("key"); err != nil {
			log.Printf("Failed to select from ClickHouse: %v", err)
		} else {
			fmt.Printf("Selected value: %s\n", value)
		}
	} else {
		fmt.Println("ClickHouse selector is not available")
	}
}
