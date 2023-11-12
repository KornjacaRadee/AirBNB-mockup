# AirBNB-mockup

Accommodation Booking System! This system is designed to facilitate the interaction between Unauthenticated Users (UU), Hosts (H), and Guests (G) in the process of discovering, managing, and booking accommodations.

## Roles in the System:

### Unauthenticated User (NK):
- Can create a new host or guest account or log in to an existing one.
- Able to browse accommodations but cannot make reservations or create new listings.

### Host (H):
- Can create and manage accommodations.
- Defines amenities, availability periods, and pricing for each accommodation.
- Can view and search all properties on their account but cannot make reservations.

### Guest (G):
- Can reserve accommodations.
- Has the ability to cancel a reservation before its start date.
- Can rate accommodations and host accounts.

## System Components:

### Client Application:
- Provides a graphical interface for users to access the system's functionalities.

### Server Application:
A microservices application that will contain the following services:

#### Auth Service:
- Manages user credentials and roles in the system.
- Handles user registration and login processes.

#### Profile Service:
- Stores basic user information such as name, gender, age, email, etc.

#### Accommodations Service:
- Manages fundamental information about accommodations (name, description, images, etc.).

#### Reservations Service:
- Controls availability periods, accommodation prices, and handles all reservations.

#### Recommendations Service:
- Supports operations for recommending accommodations to guests.

#### Notifications Service:
- Manages the storage and sending of notifications to users.