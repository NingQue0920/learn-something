### 共享锁
```sql
-- 事务1：加共享锁，读取数据
START TRANSACTION;
SELECT * FROM products WHERE id = 1 LOCK IN SHARE MODE;

-- 事务2：也能加共享锁，读取相同数据
START TRANSACTION;
SELECT * FROM products WHERE id = 1 LOCK IN SHARE MODE;

-- 事务3：尝试更新数据，发生阻塞，直到事务1或2提交或回滚
UPDATE products SET price = 100 WHERE id = 1;

-- 提交或回滚事务1和事务2后，事务3可以继续
COMMIT;

```

### 排他锁
```sql
-- 事务1：加排他锁，修改数据
START TRANSACTION;
SELECT * FROM products WHERE id = 1 FOR UPDATE;

-- 事务2：尝试读取数据，发生阻塞，直到事务1提交或回滚
SELECT * FROM products WHERE id = 1 LOCK IN SHARE MODE;

-- 事务1提交后，事务2可以继续
COMMIT;

```

### 行级锁(Row-level Locks)
```sql
-- 事务1：对ID为1的记录加行级锁
START TRANSACTION;
SELECT * FROM orders WHERE order_id = 1 FOR UPDATE;

-- 事务2：可以对ID为2的记录操作，不受影响
START TRANSACTION;
UPDATE orders SET status = 'shipped' WHERE order_id = 2;

-- 事务2对ID为1的记录操作，发生阻塞
UPDATE orders SET status = 'shipped' WHERE order_id = 1;

-- 事务1提交或回滚后，事务2可以继续
COMMIT;

```

### 表级锁(Table-level Locks)
```sql
-- 事务1：加表级锁
LOCK TABLES employees WRITE;

-- 事务2：无法对表进行任何操作，直到事务1释放锁
SELECT * FROM employees;

-- 释放表锁
UNLOCK TABLES;

```

### 意向锁(Intention Locks)
```sql
-- 事务1：加意向排他锁（IX），准备加行级排他锁
START TRANSACTION;
SELECT * FROM employees WHERE employee_id = 5 FOR UPDATE;

-- 事务2：加意向共享锁（IS），但不会与事务1冲突
SELECT * FROM employees WHERE employee_id = 6 LOCK IN SHARE MODE;

-- 事务3：尝试加表级共享锁，发生冲突，直到事务1提交
LOCK TABLES employees READ;

-- 提交事务1，释放锁
COMMIT;

```

### 间隙锁(Gap Locks)
```sql
-- 事务1：使用可重复读隔离级别，对范围加间隙锁
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
START TRANSACTION;
SELECT * FROM orders WHERE order_id BETWEEN 1 AND 10 FOR UPDATE;

-- 事务2：尝试插入一个新的记录，ID在间隙范围内，发生阻塞
INSERT INTO orders (order_id, status) VALUES (9, 'pending');

-- 事务1提交或回滚后，事务2可以继续
COMMIT;
    
```

### 死锁检测
```sql
-- 事务1：对ID为1的记录加锁
START TRANSACTION;
UPDATE products SET stock = stock - 1 WHERE id = 1;

-- 事务2：对ID为2的记录加锁
START TRANSACTION;
UPDATE products SET stock = stock - 1 WHERE id = 2;

-- 事务1：尝试对ID为2的记录加锁，发生阻塞
UPDATE products SET stock = stock - 1 WHERE id = 2;

-- 事务2：尝试对ID为1的记录加锁，产生死锁，MySQL自动回滚其中一个事务
UPDATE products SET stock = stock - 1 WHERE id = 1;

```


### 原始表结构
```sql
CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    order_id INT PRIMARY KEY AUTO_INCREMENT,
    customer_name VARCHAR(100) NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE employees (
    employee_id INT PRIMARY KEY AUTO_INCREMENT,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    department VARCHAR(50),
    hire_date DATE,
    salary DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

```