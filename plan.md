# Melo Finance Manager
This project will have an excessive number of microservices, for the purpose of learning and understanding how to build a microservice architecture, with observability, kubernetes, TDD, CI/CD, and all infrastructure provisioned through terraform, and ansible, using the AWS cloud provider.

## User Management
- Language: Python
- Framework: Flask
- Database: Redis, PostgreSQL
- ORM: SQLAlchemy
- Test: unittest
- Manual and automatic Otel instrumentation.
- Deploy: EC2
- Publishes: UserCreated, UserUpdated events (Kafka)
- Responsabilities:
  - User registration and authentication
  - Profile management
  - User roles and permissions
  - Password management

## Transaction Management
- Language: Python
- Framework: DjangoREST
- Database: MariaDB
- Test: Pytest
- ORM: Django ORM
- Deploy: ECS
- Publishes: TransactionCreated, TransactionUpdated, TransactionDeleted events (Kafka)
- Responsabilities:
  - Financial transaction creation and management
  - Transaction categorization and tagging
  - Transaction association with users and accounts
  - Recurring transactions

## Budgeting
- Language: Python
- Framework: FastAPI
- Database: MongoDB
- Test: Pytest
- Deploy: EC2
- ORM: None
- Publishes: BudgetThresholdExceeded events (Kafka)
- Responsabilities:
  - Budget creation and management for different time periods and categories
  - Budget tracking and alerts
  - Budget allocation and allocation tracking

## Reporting and Analytics
- Language: Python
- Framework: gRPC (python gRPC library)
- Database: The transaction database (PostgreSQL)
- Test: Behave
- Deploy: EC2
- ORM: SQLAlchemy
- Consumes: TransactionCreated, TransactionUpdated events (from Kafka)
- Responsabilities:
  - Financial reports generation
  - Data visualization and analysis
  - Customizable reports and dashboards
  - Calculate net worth
  - Providing insights and recommendations

## Category Management
- Language: Go
- Framework: Gin
- Database: RDS
- Test: stdlib
- Deploy: EC2
- ORM: GORM
- Responsabilities:
  - Category creation and management
  - Category tagging and categorization, and customization
  - Category hierarchy and relationships
  - Default categories

## Account Management
- Language: Go
- Framework: gRPC
- Database: DynamoDB
- Test: stdlib
- Deploy: EC2
- ORM: None
- Consumes: TransactionCreated, TransactionUpdated events (from Kafka - for balance updates)
- Responsabilities:
  - User account management ('Checking', 'Savings', etc)
  - Balance Tracking for each account

## API Gateway
- Language: Python
- Framework: FastAPI & Strawberry
- Database: None
- Test: Pytest
- ORM: None
- Communication: GraphQL
- Deploy: EC2
- Responsabilities:
  - API Gateway for all microservices


## Frontend
- Language: JavaScript
- Framework: React
- Database: None
- Deploy: S3 and CloudFront
- Responsabilities:
  - User interface for managing financial data
  - Visualization of financial data
  - Integration with the API Gateway

## Kafka
- Message broker for asynchronous communication


## Notifications
- Language: Ruby
- Framework: Ruby on Rails (API mode)
- Database: PostgreSQL (or SQLite)
- Deploy: EC2
- Test: RSpec
- Responsibilities:
  - Consumes: BudgetThresholdExceeded, UserCreated, UserUpdated events (from Kafka)
  - Manages notification preferences
  - Sends email notifications (and potentially other channels)


## CI/CD
- Github Actions

## Infrastructure Provisioning
- Terraform
- Ansible
- AWS (EC2, Elastic BeanStalk, ECS, RDS, S3, DynamoDB)

## Monitoring and Observability
- OpenTelemetry
- Grafana
- Prometheus
- Jaeger
- Loki
