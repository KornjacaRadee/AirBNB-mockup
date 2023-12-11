import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ReservationService } from '../services/reservation/reservation.service';

@Component({
  selector: 'app-accommodation-page',
  templateUrl: './accommodation-page.component.html',
  styleUrls: ['./accommodation-page.component.css']
})
export class AccommodationPageComponent implements OnInit {

  accommID = "";
  accommodation: any | null;
  availabilityPeriods: any[] = [];
  currentUserID: string = ""
  reservationData: any | null;


  constructor(private reservationService: ReservationService,private authService: AuthService, private router: Router,private route: ActivatedRoute, private accommodationService: AccomodationService){}
  ngOnInit(): void {
    if(this.authService.getUserRole() == "host"){
      this.router.navigate(['/home']);
    }
    if(!this.authService.isAuthenticated()){
      this.router.navigate(['/login']);
    }
    this.currentUserID = this.authService.getUserId();
    this.route.queryParams.subscribe(params => {
      this.accommID = params['id'];
        this.loadAccommodation(this.accommID);
        this.loadAvailability(this.accommID)

    });
  }

  isHost(): boolean{
    if(this.authService.getUserRole() == "host"){

      return true
    }else{
      return false
    }
  }
  loadAccommodation(id: string){
    this.accommodationService.getAccommodation(id).subscribe(
      (data: any[]) => {
        this.accommodation = data;
      },
      (error: any) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  loadAvailability(accommodationId: string){
    this.reservationService.getAvailability(accommodationId).subscribe(
      (data: any[]) => {
        this.availabilityPeriods = data;
      },
      (error: any) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  reservePeriod(periodid: any,start: any, end:any,): void {
      this.reservationData = {
        AvailabilityPeriodId: periodid,
        StartDate: start,
        EndDate: end,
        GuestId: this.currentUserID,
      };

    this.reservationService.postReservation(this.reservationData)
      .subscribe(
        (response) => {
          console.log('Reservation successful:', response);
          // Handle success, e.g., show a success message
        },
        (error) => {
          console.error('Reservation failed:', error);
          // Handle error, e.g., show an error message
        }
      );
  }

}
