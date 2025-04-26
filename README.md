
# Beneficiary Manager Adapter

This repository contains the backend service for the Beneficiary Manager Adapter. It includes functionality such as pagination, filtering, and dynamic data fetching for managing schemes.

## Prerequisites

Before running the project, ensure you have the following installed on your system:

- [Go](https://golang.org/dl/) 
- [PostgreSQL](https://www.postgresql.org/download/) 
- [Git](https://git-scm.com/)

## Getting Started

Follow these steps to clone the repository and run the backend service:

### 1. Clone the Repository

```bash
git clone https://github.com/ONEST-Network/Beneficiary-Manager-Adapter.git
cd Beneficiary-Manager-Adapter/backend
```
### 2. Install Dependencies

This project uses Go for the backend. You can install the dependencies by running:
```bash
go mod tidy
``` 
### 3. Set Up Environment Variables
You will need to configure some environment variables for your database connection and other settings. Create a .env file in the backend folder with the following content:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASSWORD=your_database_password
Replace your_database_name, your_database_user, and your_database_password with your PostgreSQL database credentials.
```

### 4. Run the Application
You can now run the backend server:

- Generate Go struct for the extra fields listed in the external_ref_fields.yaml.

```bash
go generate ./...
```

- Build the project using following command.

```bash
go build ./cmd/laas
```
- Run the executable.

```bash
./laas
```

### 5. Database Setup (Optional)


Let me know if you'd like any further modifications!