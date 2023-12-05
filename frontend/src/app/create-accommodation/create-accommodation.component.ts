import { Component } from '@angular/core';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { AuthService } from '../services/auth.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';


@Component({
  selector: 'app-create-accommodation',
  templateUrl: './create-accommodation.component.html',
  styleUrls: ['./create-accommodation.component.css']
})
export class CreateAccommodationComponent {
  accommodation: any = {
    name: '',
    location: '',
    minGuestNum: 0,
    maxGuestNum: 0,
    amenities: []
  };
  amenities: string[] = ['WiFi', 'Parking', 'Kitchen', 'Pool']; // Add your amenities

  constructor(
    private accommodationService: AccomodationService,
    private authService: AuthService
  ) {}

  createAccommodation() {
    if (this.authService.isAuthenticated()) {
      const token = this.authService.getAuthToken();

      const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

      this.accommodationService.createAccommodation(headers, this.accommodation).subscribe(
        (response) => {
          console.log('Accommodation created successfully', response);
        },
        (error) => {
          // Handle error, e.g., show an error message
          console.error('Failed to create accommodation', error);
        }
      );
    } else {
      // User is not authenticated, handle accordingly
      console.error('User is not authenticated');
    }
  }
}
