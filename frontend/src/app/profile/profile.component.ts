import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ProfilesService } from '../services/profile/profiles.service';
import { CommonModule } from '@angular/common';
import { ReservationService } from '../services/reservation/reservation.service';
import { NgToastService } from 'ng-angular-popup';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  user: any;
  id: string = '';
  accomms: any[] = [];
  reservations: any[] = [];
  showAccommodations = false;
  showReservations = false;

  profile: any;

  constructor(
    private toastr: ToastrService,
    private toast: NgToastService,
    private reservationService: ReservationService,
    private profileService: ProfilesService,
    private authService: AuthService,
    private accommodationService: AccomodationService,
    private router: Router
  ) {}

  ngOnInit(): void {
    if (!this.authService.isAuthenticated()) {
      this.router.navigate(['/login']);
    }
    this.accomms = [];
    //this.getUserAccommodations()
    this.loadUserDetails();
    console.log(this.accomms);
  }
  tempLoadAccoms(): void {
    this.accommodationService.getAccomodations().subscribe(
      (data: any[]) => {
        this.accomms = data;
        this.accomms = data.filter(
          (accommodation) => accommodation.owner.id === this.id

        );
        console.log(this.accomms);
        if(this.accomms.length == 0){
          this.toastr.error('User has no accommodations');
        }
      },
      (error: any) => {
        this.toastr.error('Error loading accommodations');
        console.error('Error fetching accommodations:', error);
      }
    );
  }
  deleteAccomm(id: string): void{
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    this.accommodationService.deleteAccommodation(id,headers).subscribe(
      (response) => {
        this.tempLoadAccoms();
        this.toastr.success("Successfully deleted accommodation");
      },
      (error: any) => {
        this.toastr.error('Error deleting accommodation');
        console.error('Error deleting:', error);
      }
    );
  }

  getUserAccommodations(): void {
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

    this.accommodationService.getUserAccommodations(headers).subscribe(
      (data: any[]) => {
        console.log(data);
        this.accomms = data;
      },
      (error: any) => {
        this.toastr.error('User has no accommodations');
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  getUserReservations(): void {
    this.reservationService.getUserReservations(this.id).subscribe(
      (data: any[]) => {
        console.log(data);
        this.reservations = data;
        console.log(this.reservations)
        if(this.reservations == null){
          this.toastr.error('User has no reservations');
        }
      },
      (error: any) => {
        this.toastr.error('User has no reservations');
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  logout(): void {
    this.authService.logout();
  }

  isHost(): boolean {
    if (this.authService.getUserRole() == 'host') {
      return true;
    } else {
      return false;
    }
  }

  deleteProfile() {
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

    this.authService.deleteUser(headers).subscribe(
      (response) => {
        this.toastr.success('Profile and account deleted');
        this.authService.logout();
      },
      (error) => {
        // Handle error, e.g., show an error message
        this.toastr.error('Failed to delete user');
        console.error('Failed to delete user', error);
      }
    );
  }

  loadUserDetails() {
    this.id = this.authService.getUserId();
    this.authService.getUserById(this.id).subscribe(
      (response) => {
        // Map the response to the 'user' property
        this.user = response;
        this.loadProfileDetails();
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }
  toggleAccommodations() {
    this.tempLoadAccoms();
    this.showAccommodations = !this.showAccommodations;
  }

  toggleReservations() {
    this.getUserReservations();
    this.showReservations = !this.showReservations;
  }

  loadProfileDetails() {
    this.profileService.getProfileByEmail(this.user.email).subscribe(
      (response) => {
        // Map the response to the 'user' property
        this.profile = response;
        console.log(this.profile);
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }
}
