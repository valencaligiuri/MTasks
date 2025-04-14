# MTasks

MTasks is a web application designed to help you efficiently manage your tasks.

## Features

- **Create, Read, Update, and Delete (CRUD) Tasks:** Manage your tasks with ease.
- **Dark/Light Theme:** Toggle between light and dark themes for a personalized experience.
- **Progress Tracking:** Monitor the completion status of your tasks.
- **User-Friendly Interface:** Enjoy a simple and intuitive user experience.
- **Authentication:** Secure access to your tasks with password-protected login.

## Technologies Used

- **Frontend:** Vue 3, TypeScript, TailwindCSS, Vite
- **Backend:** Go (Gin Framework)
- **Database:** SQLite

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/valencaligiuri/MTasks.git
    ```

2. Navigate to the project directory:

    ```bash
    cd MTasks
    ```

3. Install backend dependencies:

    ```bash
    go mod tidy
    ```

4. Navigate to the frontend directory and install dependencies:

    ```bash
    cd frontend
    npm install
    ```

5. Configure the database:
    - The backend uses SQLite. The database file (`tasks.db`) will be created automatically when the server starts.

## Usage

1. Start the backend server:

    ```bash
    go run main.go
    ```

    - You will be prompted to set a password for the server.

2. Start the frontend development server:

    ```bash
    npm run dev
    ```

3. Open your web browser and navigate to `http://localhost:8080`.

## Contributing

Contributions are welcome! Fork the repository and submit a pull request with your changes. Please follow these guidelines:

- Ensure your code adheres to the project's coding standards.
- Write clear and concise commit messages.
- Test your changes thoroughly.

## License

This project is licensed under the MIT License.

## Contact

For any inquiries, please contact me at [valencaligiuri06@gmail.com](mailto:valencaligiuri06@gmail.com).