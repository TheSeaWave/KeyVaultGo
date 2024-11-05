In-Memory Key-Value Store with HTTP Interface
Overview
This project is a lightweight in-memory key-value database that provides an HTTP interface for simple and fast storage and retrieval operations. It supports persistence to disk for data durability and includes a feature for automatically removing expired records.

Features
In-Memory Storage: Fast data access due to in-memory operations.
HTTP API: Easy-to-use HTTP endpoints for CRUD operations.
Persistent Storage: Optional persistent storage to retain data across server restarts.
Automatic Expiration: Configurable TTL (Time-To-Live) for automatic deletion of outdated records.


запуск бенчмарков 
go test -bench=. -benchtime=1s -benchmem > /home/student/dz1/benchmark_results/benchmark_results.txt