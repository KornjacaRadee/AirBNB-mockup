import { Component } from '@angular/core';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { AuthService } from '../services/auth.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NgToastService } from 'ng-angular-popup';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';


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
    private toastr: ToastrService,
    private toast: NgToastService,
    private router: Router,
    private accommodationService: AccomodationService,
    private authService: AuthService,
  ) {}

  createAccommodation() {
    this.showSuccess();
    if (this.authService.isAuthenticated()) {
      const token = this.authService.getAuthToken();

      const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

      this.accommodationService.createAccommodation(headers, this.accommodation).subscribe(
        (response) => {
          this.toastr.success('Accommodation created successfully! Add pictures on your profile page! :D');
          console.log('Accommodation created successfully', response);
        },
        (error) => {
          // Handle error, e.g., show an error message
          this.toastr.error('Failed to create accommodation!');
          console.error('Failed to create accommodation', error);
        }
      );
    } else {
      // User is not authenticated, handle accordingly
      this.toastr.error('Please log in');
      console.error('User is not authenticated');
    }
  }



  showSuccess() {
    this.toast.success({detail:"SUCCESS",summary:'Your Success Message',duration:5000});
  }
}
