#Remove the previous build
rm bootstrap
rm transaction.zip
# Copying the env vars file for AWS
cp -r ../../files .
#Build the lambda
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
#Compress the executable
zip -qr transaction.zip bootstrap app.env files
# Declutter folder
rm -rf files