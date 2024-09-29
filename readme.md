Certainly! Below is a **rewritten `README.md`** for your **RTO Attendance Tracker** project. This updated version reflects the current state of your application, incorporating the features you've implemented so far, best practices, and enhanced clarity to help users and contributors understand and engage with your project effectively.

---

# RTO Attendance Tracker

A comprehensive web application to track and visualize office attendance, enabling efficient planning and management of in-office and remote workdays. Designed to support Return to Office (RTO) policies, this tool helps users balance their work arrangements seamlessly.

![RTO Attendance Tracker Screenshot](path-to-screenshot.png) <!-- Replace with an actual screenshot -->

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Demo](#demo)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Technology Stack](#technology-stack)
- [Contributing](#contributing)
- [Testing](#testing)
- [Roadmap](#roadmap)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Introduction

With the evolving landscape of work environments, managing in-office and remote days has become essential for both employees and employers. Companies are instituting RTO policies that require a certain number of days in the office, and keeping track of these schedules can be challenging. The **RTO Attendance Tracker** addresses this need by providing a personal dashboard to monitor and visualize attendance, incorporating holidays and planned vacations to offer a clear overview of your work arrangements.

## Features

- **Attendance Preferences**
  - Set default in-office and remote days for each week.
  - Easily update preferences at any time.
  
- **Holiday and Vacation Import**
  - Load holidays and planned vacations from a JSON file.
  - Automatically integrates these events into your attendance schedule.
  
- **Visual Calendar**
  - Interactive calendar displaying all days within the quarter.
  - Color-coded indicators:
    - **Green** for In-Office Days
    - **Red** for Remote Days
    - **Blue** for Holidays
    - **Purple** for Vacations
    - **Gray** for Weekends
    
- **Dynamic Updates**
  - Real-time updates to the calendar when preferences are modified.
  
- **Data Persistence**
  - Saves user preferences and events locally using JSON files.
  
- **Statistics and Analytics**
  - Calculates and displays average in-office days per week.
  - Visual bar graphs representing attendance trends.
  
- **Structured Logging**
  - Integrated with `slog` for efficient and structured logging.
  
- **Responsive Design**
  - Optimized for various devices, ensuring a seamless experience on desktops, tablets, and mobile phones.
  
- **User-Friendly Interface**
  - Intuitive navigation and interactive elements for easy usage.

## Demo

![Calendar View](path-to-calendar-screenshot.png) <!-- Replace with an actual screenshot -->
*An overview of the calendar displaying in-office and remote days, along with holidays and vacations.*

## Getting Started

Follow these instructions to set up and run the project on your local machine.

### Prerequisites

- **Go**: Version 1.20 or higher
- **Git**: For cloning the repository
- **Node.js & npm**: If you plan to manage frontend dependencies separately
- **Docker**: Optional, for containerized deployment

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/rto-attendance-tracker.git
   ```

2. **Navigate to the Project Directory**

   ```bash
   cd rto-attendance-tracker
   ```

3. **Install Dependencies**

   ```bash
   go mod download
   ```

4. **Set Up Configuration Files**

   Ensure the `data` directory contains the necessary JSON files:

   - `holidays.json`: List of holidays and vacations.
   - `events.json`: Attendance events.
   - `preferences.json`: User preferences.

   Example structure for `holidays.json`:

   ```json
   [
     {
       "date": "2024-11-28",
       "name": "Thanksgiving",
       "type": "holiday"
     },
     {
       "date": "2024-12-25",
       "name": "Christmas",
       "type": "holiday"
     }
   ]
   ```

5. **Run the Application**

   ```bash
   go run cmd/main/main.go
   ```

   The application will start on `http://localhost:8761`.

### Docker Deployment (Optional)

1. **Build the Docker Image**

   ```bash
   docker build -t rto-attendance-tracker .
   ```

2. **Run the Docker Container**

   ```bash
   docker run -d -p 8761:8761 --name rto-attendance-tracker rto-attendance-tracker
   ```

   Access the application at `http://localhost:8761`.

## Usage

### Setting Attendance Preferences

1. **Navigate to Preferences**

   Click on the **Prefs** button or go to `/prefs`.

2. **Configure Default Days**

   - Enter the default in-office days using abbreviations (e.g., `M,T,W,Th,F`).
   - Set the target number of in-office days per week (e.g., `2.5`).

3. **Save Preferences**

   Click the **Save Preferences** button to apply changes.

### Importing Holidays and Vacations

1. **Prepare `holidays.json`**

   Ensure your JSON file follows the required structure with `date`, `name`, and `type` fields.

2. **Load the File**

   The application automatically reads the `holidays.json` file on startup. To update, replace the existing file and restart the application.

### Viewing and Managing Events

1. **Access the Calendar**

   Visit the home page (`/`) to view your attendance calendar.

2. **Toggle Attendance Status**

   - Click on any **In Office** or **Remote** indicator within the calendar to toggle the status.
   - The system will update the event and recalculate statistics accordingly.

3. **Add New Events**

   - Navigate to `/add-event` or click the **Add Event** button.
   - Fill in the event details and submit the form to add holidays, vacations, or specific attendance events.

4. **View All Events**

   - Go to `/events` or click the **Events** button to see a comprehensive list of all events.

## Technology Stack

- **Backend:**
  - [Go](https://golang.org/) with the [Echo](https://echo.labstack.com/) framework
  - [slog](https://pkg.go.dev/log/slog) for structured logging
  - [GORM](https://gorm.io/) for ORM (if integrated in future steps)

- **Frontend:**
  - HTML5 & CSS3 for structure and styling
  - [Font Awesome](https://fontawesome.com/) for icons
  - [Chart.js](https://www.chartjs.org/) for data visualization
  - [jQuery](https://jquery.com/) for DOM manipulation and AJAX (optional)

- **Deployment:**
  - Docker for containerization (optional)

## Contributing

Contributions are welcome! Follow these steps to contribute to the project:

1. **Fork the Repository**

   Click the **Fork** button at the top-right corner of the repository page.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/yourusername/rto-attendance-tracker.git
   ```

3. **Create a New Branch**

   ```bash
   git checkout -b feature/YourFeatureName
   ```

4. **Commit Your Changes**

   ```bash
   git commit -m "Add Your Feature"
   ```

5. **Push to Your Fork**

   ```bash
   git push origin feature/YourFeatureName
   ``` 
 