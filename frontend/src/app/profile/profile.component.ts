import { Component, OnInit, Output, EventEmitter,Input  } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import {  HttpHeaders } from '@angular/common/http';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ProfilesService } from '../services/profile/profiles.service';
import { ReservationService } from '../services/reservation/reservation.service';
import { NgToastService } from 'ng-angular-popup';
import { ToastrService } from 'ngx-toastr';
interface Review {
  id: number;
  GuestId: string;
  HostId: string;
  Time: string;
  AccommodationId: string;
  Guest: any,
  Rating: number;
}

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css'],
})
export class ProfileComponent implements OnInit {
  @Input() reviews: Review[] = [];

  @Input() accomReviews: Review[] = [];
  user: any;
  id: string = '';
  accomms: any[] = [];
  notifications: any[] = [];
  reservations: any[] = [];
  showAccommodations = false;
  showReservations = false;
  klasica = false;
  showNotifications: boolean = false;
  showRatings: boolean = false;
  showOverlay: boolean = false;
  reservationHistory: any[] = [];
  activeReservations: any[] = [];
  stars: number[] = [1, 2, 3, 4, 5];
  currentRating: number = 0;
  pictures: any[] = [];
  picsdata: any[] = [];


  @Output() ratingChanged: EventEmitter<number> = new EventEmitter<number>();


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
    this.loadNotifications();

  }
  tempLoadAccoms(): void {
    this.accommodationService.getAccomodations().subscribe(
      (data: any[]) => {
        this.accomms = data;

        this.accomms = data.filter(
          (accommodation) => accommodation.owner.id === this.user.id

        );
        this.accomms.forEach(ac =>{
          this.loadPictures(ac.id)
          ac.picture = this.picsdata;

        })

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
  toggleOverlay(): void {
    this.reservationHistory = [];
    this.getUserReservations();
    this.showOverlay = !this.showOverlay;
  }


  toggleRatings(): void {
    this.loadUserRatings();
    this.loadAccommodationRatings();
    this.showRatings = !this.showRatings;
    this.showAccommodations = false;
  }

  loadNotifications(): void{
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    this.profileService.getUserNotifications(headers).subscribe(
      (data: any[]) => {
        this.notifications = data;

        data.forEach(not =>{
          not.time = new Date(not.time).toISOString().split('T')[0];
        } )
      },
      (error: any) => {
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

        this.accomms = data;
      },
      (error: any) => {
        this.toastr.error('User has no accommodations');
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  getUserReservations(): void {
    this.activeReservations = [];
    this.reservationHistory = [];
    this.reservations = [];
    this.reservationService.getUserReservations(this.id).subscribe(
      (data: any[]) => {

        this.reservations = data;

        if(this.reservations == null){

        }else{
          data.forEach(not =>{
            var todaysDate = new Date().toISOString().split('T')[0]
            not.StartDate = new Date(not.StartDate).toISOString().split('T')[0];
            not.EndDate = new Date(not.EndDate).toISOString().split('T')[0];
            console.log(not.EndDate)
            console.log(todaysDate)
            this.loadPictures(not.AccommodationId);
            if (not.EndDate < todaysDate){
              this.accommodationService.getAccommodation(not.AccommodationId).subscribe(
                (data: any[]) => {
                  not.accommodation = data;
                  not.picture = this.picsdata;
                  this.reservationHistory.push(not)


                },
                (error: any) => {
                  console.error('Error fetching accommodations:', error);
                }
                );
            }else{
              this.accommodationService.getAccommodation(not.AccommodationId).subscribe(
                (data: any[]) => {
                  not.accommodation = data;
                  not.picture = this.picsdata;
                  this.activeReservations.push(not)


                },
                (error: any) => {
                  this.toastr.error('User has no reservations');

                  console.error('Error fetching accommodations:', error);
                }
              );

            }
          } )

        }
      },
      (error: any) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  fill(star: number): void {
    this.currentRating = star;
  }

  reset(): void {
    if (this.currentRating === 0) {
      return;
    }
    this.currentRating = 0;
  }

  rate(star: number,reservation: any): void {

    this.currentRating = star;
    this.ratingChanged.emit(this.currentRating);
    const reservationJSON = {
      HostId: reservation.accommodation.owner.id,
      GuestId: reservation.GuestId,
      AccommodationId: reservation.accommodation.id,
      time: new Date(),
      rating: star
    }
    // reservationJSON = JSON.stringify(reservationJSON);
    this.accommodationService.rateAccommodation(reservationJSON).subscribe(
      (data: any[]) => {
        this.toastr.error('Successfully rated accommodation.');
      },
      (response: any) => {
        if(response.error && response.error.includes('502')){
          this.toastr.success('Successfully rated accommodation.');
        }else{
          this.toastr.error('Failed to rate accommodation.');
        }
        this.toastr.error('Error with rating. Try again later! ');
      }
    );

  }
  loadPictures(id: string){
    this.pictures = [];
    this.picsdata = [];
    this.accommodationService.getAccommodationPictures(id).subscribe(
      (data: any[]) => {
        this.pictures = data;
        this.pictures.forEach(pic =>{
          this.picsdata.push(pic.data)
        })
      },
      (error: any) => {
        this.toastr.error('Error fetching accommodation!');
        console.error('Error fetching accommodations:', error);
      }
    );
  }


  rateHost(star: number,reservation: any): void {

    this.currentRating = star;
    this.ratingChanged.emit(this.currentRating);
    const reservationJSON = {
      HostId: reservation.accommodation.owner.id,
      GuestId: reservation.GuestId,
      time: new Date(),
      rating: star
    }
    // reservationJSON = JSON.stringify(reservationJSON);
    this.profileService.rateHost(reservationJSON).subscribe(
      (data: any[]) => {
        this.toastr.success('Successfully rated accommodation.');
      },
      (response: any) => {
        if(response.error && response.error.includes('502')){
          this.toastr.success('Successfully rated accommodation.');
        }else{
          this.toastr.error('Failed to rate accommodation.');
        }
        this.toastr.error('Error with rating. Try again later! ');
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
    this.showRatings = false;
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
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }

  loadUserRatings(){
    this.profileService.getHostRatings(this.user.id).subscribe(
      (data) => {
        this.reviews = data;
        this.reviews.forEach(ac => {
          ac.Time =  new Date(ac.Time).toISOString().split('T')[0];
          this.authService.getUserById(ac.GuestId).subscribe(
            (response) => {
              ac.Guest = response;
            },
            (error) => {
              console.error('Error fetching user details', error);
            }
          );
        })

      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }


  loadAccommodationRatings(){
    this.profileService.getAccommodationRatings(this.user.id).subscribe(
      (data) => {
        this.accomReviews = data;
        this.accomReviews.forEach(ac => {
          ac.Time =  new Date(ac.Time).toISOString().split('T')[0];
          this.authService.getUserById(ac.GuestId).subscribe(
            (response) => {
              ac.Guest = response;
            },
            (error) => {
              console.error('Error fetching user details', error);
            }
          );
        })

      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }

}
