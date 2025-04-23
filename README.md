# Lilo Finance Manager

[![Work in Progress](https://img.shields.io/badge/Status-Work%20in%20Progress-yellow.svg)](https://github.com/your-username/lilo-finance-manager)

## Overview

Lilo Finance Manager is a personal finance web application built with a microservices architecture to help users manage their finances effectively. This project aims to provide a comprehensive suite of tools for tracking transactions, managing accounts and categories, creating budgets, setting up notifications, and generating insightful reports.

**Important Note:** Currently, the **User Management** service is the only microservice with its initial version (v1) developed. The other services are either in the planning stages or under active development. We also plan to integrate Kafka messaging for inter-service communication in the future.

## Microservices Architecture

The application is composed of the following microservices:

* [Work in Progress - v1 Ready] **User Management:** A REST API (Python, Flask, PostgreSQL, SQLAlchemy) for handling user registration, authentication, and profile management.
* [Not started] **Transaction Management:** A REST API (Python, Django REST Framework, MariaDB) for recording and managing financial transactions.
* [Not started] **Account Management:** A gRPC API (Go, gRPC-gateway, DynamoDB) for managing user accounts and balances.
* [Not started] **Category Management:** A REST API (Go, Gin, potentially MSSQL) for organizing transactions into categories.
* [Not started] **Budgeting Management:** A REST API (Python, FastAPI, MongoDB) for creating and tracking budgets.
* [Not started] **Notification Management:** (Technology to be determined) for sending users relevant financial notifications.
* [Not started] **Reporting and Analytics:** A gRPC API (Python) for generating financial reports and insights.
* [Not started] **API Gateway:** A GraphQL API (Python, FastAPI, Strawberry) serving as the entry point for the frontend and orchestrating requests to backend services.
* [Not started] **Frontend:** (TypeScript, React) - The user interface for interacting with the application.

## Development Setup (Dev Environment)

The development environment is set up using:

* **Operating System:** Ubuntu running in a Multipass VM (provisioned via Terraform).
* **Configuration Management:** Ansible.
* **Service Mesh:** Istio.
* **Observability Stack:** OpenTelemetry Collector, Prometheus, Grafana, Loki, Jaeger.

## Production Deployment (Prod Environment)

The production environment will be deployed on AWS using a combination of services:

* ECS (Elastic Container Service)
* Elastic Beanstalk
* RDS (Relational Database Service)
* S3 (Simple Storage Service)
* EBS (Elastic Block Store)
* DynamoDB
* EC2 (Elastic Compute Cloud)

## CI/CD

Continuous Integration and Continuous Deployment pipelines will be automated using GitHub Actions to ensure efficient and reliable software delivery.

## Development Workflow

The project follows the Gitflow branching model:

* `main`: Represents the production-ready state.
* `develop`: The main integration branch for ongoing development.
* `feature/*`: Branches for developing specific features.
* `release/*`: Branches for preparing a new release.
* `hotfix/*`: Branches for addressing critical issues in production.

## Future Enhancements

* Integration of Kafka messaging for asynchronous communication between microservices.
* Development and deployment of the remaining microservices.
* Implementation of comprehensive test suites for all services.
* Design and development of the frontend user interface.
* Setting up robust monitoring and alerting in both development and production environments.

## Getting Started (For Developers)

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/lilo-finance-manager.git](https://github.com/your-username/lilo-finance-manager.git)
    cd lilo-finance-manager
    ```

2.  **Set up the development environment:** Ensure you have make Terraform, Ansible and Multipass installed. Navigate to the `deploy/infrastructure` directory and run:
    ```bash
    make dev
    ```

3.  **Explore the User Management service:** Navigate to the `user-management` directory for specific instructions on running and testing this service.

## Contributing

Contributions are welcome! Please follow the Gitflow workflow and submit pull requests for review.

## License

GPL-3.0 License.

## Contact

[Edmilson Rodrigues](mailto:edmilson.monteiro.rodrigues@gmail.com)
