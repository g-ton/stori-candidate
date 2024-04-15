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
 - S3
 - API Gateway
 - Aurora RDS

# Consuming API endpoints

Here the Postman collection to consume the endpoints locally through localhost:8282 

### sendSummaryInfoByFile
The easiest endpoint to consume is `sendSummaryInfoByFile` because it's not necessary to create previously an account nor a transaction, this endpoint can read two possible files **./files/txns.csv**, **./files/txns2.csv**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/552546e6-bd93-4e0f-868c-350c45395618)

**Input:**

./files/txns.csv

![image](https://github.com/g-ton/stori-candidate/assets/13384146/208c19af-b4a0-4272-9a8d-8a7d11a9d32c)


**Expected outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/758f4622-d463-4ae4-aafc-17c15bbf7593)

**Got outcome:**

An email sent to the customer

![image](https://github.com/g-ton/stori-candidate/assets/13384146/1591cb61-ba9d-4776-aec8-1d6d5c390940)

### createAccount

**Input:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/f61719a6-82b3-4ef8-8ad3-474302f9e180)

**Outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/79c6585b-0f83-4b11-9720-a11c565043d8)

### getAccount

**Input:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/a9821fbe-ff68-4aa7-b728-8a533a7af6df)

**Outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/2023afb2-815a-4050-913e-8dfecb4537a6)

### createTransaction

**Input:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/c151c769-130c-4657-b190-71dfc5e7ab58)

![image](https://github.com/g-ton/stori-candidate/assets/13384146/50b898d1-8c90-472f-ba55-d670624db9a6)


**Outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/4f6c8181-affe-4294-ac8a-2346d939def5)

![image](https://github.com/g-ton/stori-candidate/assets/13384146/8bda5977-89ce-49a5-bbd0-5d25f360f716)



### getTransaction

**Input:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/df639a39-d10f-4e2b-83a5-22265abb75be)

![image](https://github.com/g-ton/stori-candidate/assets/13384146/6d838394-84db-4f5e-bf37-54ae343576d0)


**Outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/9372d475-fd95-4b8f-b7c4-391e18998be8)

![image](https://github.com/g-ton/stori-candidate/assets/13384146/3ebde881-2ecd-41e4-b5f0-a3f426cf1b84)


### sendSummaryInfoByDB

**Input:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/7be3d772-03e1-47f9-a175-fbdb95fad565)


**Expected outcome:**

![image](https://github.com/g-ton/stori-candidate/assets/13384146/5c090cef-e315-4af0-b825-99d9998df601)


**Got outcome:**

An email sent to the customer

![image](https://github.com/g-ton/stori-candidate/assets/13384146/b88230db-8155-4732-917d-a3ddec986e61)


















