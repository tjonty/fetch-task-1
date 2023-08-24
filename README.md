# fetch-task-1
This is the first task from Fetch: `Receipt Processor`
- https://github.com/fetch-rewards/receipt-processor-challenge

## How to Run

1. Clone the repository: `git clone https://github.com/tjonty/fetch-task-1.git`
2. Navigate to the project directory: `cd fetch-task-1`
3. Build Docker Image: `docker build -t image-name .`
4. Run Docker Image: `docker run -p 8080:8080 image-name`

## Testing

Using Postman
1. Send POST Request:

- URL: `http://localhost:8080/receipts/process`
- Method: `POST`
- Request Body: `Raw JSON payload`
- Response: `ID` (generated by code)

2. Send GET Request:

- URL: `http://localhost:8080/receipts/:id/points`
- Method: `GET`
- URL Parameters: `Replace: id` with the ID received from the POST request
- Response: `Result` (number of points)