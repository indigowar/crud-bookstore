# Book-Store

Task: Implement a simple web server that exposes an API endpoint to manage a collection of books.

Requirements:

- The server should expose the following API endpoints:
    - `/books` (GET): Retrieve a list of all books.
    - `/books/{id}` (GET): Retrieve details of a specific book by its ID.
    - `/books` (POST): Add a new book to the collection.
    - `/books/{id}` (PUT): Update the details of a specific book by its ID.
    - `/books/{id}` (DELETE): Delete a specific book by its ID.
- The book structure should include attributes like title, author, publication year, and ISBN.
- The server should store the books in memory (no need for a persistent database).
- Proper error handling should be implemented, including appropriate status codes for different scenarios.
- Implement basic validation for incoming data to ensure required fields are present, and data is in the correct format.
- Include appropriate documentation (e.g., Swagger/OpenAPI) for the API endpoints.

Additional considerations:

- The code should follow Go best practices, including proper structuring and modularity.
- Use appropriate Go packages and libraries to simplify the implementation.
- Include unit tests to ensure the correctness of the implemented functionality.
- Demonstrate understanding of RESTful API principles and good API design practices.