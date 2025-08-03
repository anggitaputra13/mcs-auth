# mcs-auth
#microservices auth using go gin postgresql redis jwt by anggita putra

#How to Run the Application with Docker
#Build and start the services:
docker-compose up --build

#To run in detached mode:
docker-compose up -d --build

#To stop the services:
docker-compose down

#To view logs:
docker-compose logs -f

#Accessing the Services
Authentication Service: http://localhost:8000
Swagger UI: http://localhost:8000/swagger/index.html
PostgreSQL: Accessible on localhost:5432 (username: postgres, password: postgres)
Redis: Accessible on localhost:6379