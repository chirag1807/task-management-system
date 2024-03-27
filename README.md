# task-management-system

Task management system is developed in Golang. It allows users to register, login, create teams, create tasks, and assign tasks to other users or teams. Users can also change task deadlines and statuses.

# Features
- User Registration: Users can create an account to access the task manager system.
- User Authentication: Secure login functionality to authenticate users.
- Team Creation: Users can create teams and manage team members.
- Task Creation: Users can create tasks with details like title, description, and deadline.
- Task Assignment: Tasks can be assigned to individual users or entire teams.
- Deadline Management: Assignees can change the deadline of tasks as required.
- Task Status Update: Assignees can update the status of tasks (e.g., In Progress, Completed).

# Tech Stack 💻
- GO 1.22
- CockroachDB
- Redis
- RabbitMQ

## Installation

### Using Docker:

1. Clone the repository using git clone:
```
$ git clone https://github.com/chirag1807/task-management-system
$ cd task-management-system
```
2. Create .config directory to in current directory and copy the `.env.example` file to new `.config/.env` file:
```
$ mkdir .config
$ cp .env.example .config/.env
```
3. Spin up the docker container:
```
$ docker-compose up
```
If permission error occurs, run command as root:
```
$ sudo docker-compose up
```
- The server will start listening on port `9090`

### Using Source Code:

#### Prerequisites you need to set up on your local computer:
1. [Golang](https://go.dev/doc/install)
2. [Redis](https://redis.io/download/)
3. [Cockroach](https://www.cockroachlabs.com/docs/releases/)
4. [RabbitMQ](https://www.rabbitmq.com/download.html)
5. [Dbmate](https://github.com/amacneil/dbmate#installation)

#### Getting Started:

1. Clone the repository using git clone:
```
1) git clone https://github.com/chirag1807/task-management-system
2) cd task-management-system
```
2. Create .config directory to in current directory and copy the `.env.example` file to new `.config/.env` file:
```
$ mkdir .config
$ cp .env.example .config/.env
```
3. Create `.env` file in current directory and update below configurations:
   1. Add Cockroach database URL in `DATABASE_URL` variable.
4. Run `dbmate migrate` to migrate database schema.
5. Run `go mod vendor` to install all the dependencies.
6. Run `go run cmd/main.go` to run the programme.

## API Documentation:

After executing run command, open your favorite browser and type below URL to open API documentation.
```
http://localhost:9090/swagger/index.html/
```

## Contribution
[<img alt="Chirag Makwana" src="https://github.com/chirag1807/task-management-system/assets/94277910/0e27ad00-c278-4eea-81df-8c3096c1ed2c" width="84" height="100" style="border-radius: 50%;" />](https://github.com/chirag1807)