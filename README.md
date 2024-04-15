# stori-candidate
A repo with the challenge for Stori Software Engineer position

The code stored in this repo has been tested (Applying units tests), the arrow icon at the top of the repo means that all tests passed successfully

![image](https://github.com/g-ton/stori-candidate/assets/13384146/da246972-96e1-4855-9fc2-ab8d6143891e)

Here the coverage for api and db modules:

![image](https://github.com/g-ton/stori-candidate/assets/13384146/77fec644-42f1-4a37-8c81-dd51f6d0c966)

# Instructions to run the project locally


**NOTE**: It's necessary first to have installed docker in your machine

 1. Clone this repo in your computer
 2. Go to the root folder of the project: `cd stori-candidate`
 3. In the terminal type this command: `docker compose up`
 4. The step number three will create two docker images (One for the stori project and other for postgres)
 5. If all goes well, it will be possible to see the log with the init for gin
    ![image](https://github.com/g-ton/stori-candidate/assets/13384146/482b5cc9-e0c8-43ed-8fe5-9c5a4978e7cc)

 7. Also if you open other terminal and type this: `docker ps`, it will be possible to see the containers running for the stori project (NOTE: Don't close the terminal where you ran the docker compose up command, otherwise, the app will be terminated)
    ![image](https://github.com/g-ton/stori-candidate/assets/13384146/915c678f-a5e6-48fa-9dc6-f3dca2eb88e4)

 8. The app is running on **localhost:8282**, the all endpoints will be explained in detail in next steps in this README, however, You can take a look at the swagger documentation available typing this URL in your browser: **http://localhost:8282/swagger/index.html**, there, it will be possible to see the endpoints for the system and their parameters and descriptions
    ![image](https://github.com/g-ton/stori-candidate/assets/13384146/c346be48-f27f-4b52-a46a-a8dd81215086)



# Brief technical description of the system

## Database
It's used a **postgres** db in order to store the account and transaction entities, you can find the diagram inside the project in **files/DB_Diagram.pdf** and it looks like this:

![image](https://github.com/g-ton/stori-candidate/assets/13384146/26c3d637-dcdd-4f0a-a995-40c03f8f3ad1)

An account is necessary to be avialbe to create a transaction, that means We need to create an account before and then with the ID of that account we can create one or more transactions

Fields for Account:

![image](https://github.com/g-ton/stori-candidate/assets/13384146/10327c68-67d8-4c24-bebe-c5a2ee365b29)


Fields for Transaction:

![image](https://github.com/g-ton/stori-candidate/assets/13384146/5c78ee8c-52a8-4413-9ae7-6bc58ed2af0b)

## Technology stack

This is the technology stack used:

**Golang**

 - Golang 1.21
 - Gin gonic as web framework
 - Testify and faker for unit tests
 - Sqlc to generate DB models
 - Mockgen to generate mocks interfaces for email and db
 - Swaggo to launch swagger documentation
 - Viper to handle config env vars
 
 **General**
 
 - Docker 20.10.12
 - Lambda functions AWS







