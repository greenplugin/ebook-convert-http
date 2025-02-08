
# **Repository Description**
This repository contains the codebase for a **Go application**.
The application is built to serve a specific functionality.

Below are the instructions to **run the app in a Docker container**.

# **How to Run the Application in Docker**

1. **Build the Docker image:**
    ```bash
    docker build -t <image-name>:<tag> .
    ```
2. **Run the Docker container:**
    ```bash
    docker run -d -p <host-port>:<container-port> <image-name>:<tag>
    ```


# **API Endpoints**
Below are the details of the available API endpoints in the application.
### **API Endpoints**

#### 1. **Get a list of recipes**
- **Endpoint**: `/recipes`
- **Method**: `GET`
- **Description**: Returns a list of recipes as a string.
- **Responses**:
   - **200 OK**:
     ```json
     [
       "Spaghetti Carbonara",
       "Chicken Curry",
       "Beef Stroganoff"
     ]
     ```
   - **500 Internal Server Error**:
     ```json
     {
       "error": "Unable to retrieve recipes"
     }
     ```

#### 2. **Check server health**
- **Endpoint**: `/health`
- **Method**: `GET`
- **Description**: Checks the server's health.
- **Responses**:
   - **200 OK**:
     ```
     OK
     ```
   - **500 Internal Server Error**:
     ```
     UNHEALTHY
     ```

#### 3. **Convert files to another format**
- **Endpoint**: `/convert`
- **Method**: `POST`
- **Description**: Converts a file to a specified format using the `ebook-convert` utility.
- **Request Body**:
   - **Content-Type**: `multipart/form-data`
   - **Parameters**:
      - `file` (required): The file to be uploaded and converted.
      - `convert-to` (required): The target format for conversion (e.g., `pdf`, `mobi`, `epub`).
      - `convert-options` (optional): Additional options for the `ebook-convert` utility in the format `"key=value"`, separated by spaces (e.g., `margin-top=10 margin-right=5`).
- **Responses**:
   - **200 OK**:
      - The file is successfully converted. The binary content of the file is returned.
      - **Headers**:
         - `Content-Disposition`: Specifies the name of the converted file, e.g.:
           ```
           attachment; filename=converted-file.pdf
           ```
   - **400 Bad Request**:
     ```json
     {
       "error": "Target format not specified (convert-to parameter)"
     }
     ```
   - **405 Method Not Allowed**:
     ```json
     {
       "error": "Only POST requests are supported"
     }
     ```
   - **500 Internal Server Error**:
     ```json
     {
       "error": "Error converting file: <details>"
     }
     ```