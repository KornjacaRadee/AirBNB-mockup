import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ReservationService } from '../services/reservation/reservation.service';
import { DatePipe } from '@angular/common';
import { formatDate } from '@angular/common';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-accommodation-page',
  templateUrl: './accommodation-page.component.html',
  styleUrls: ['./accommodation-page.component.css']
})
export class AccommodationPageComponent implements OnInit {
  selectedStartDate: string | null = "";
  selectedEndDate: string | null = "";
  parsedStartDate: any;
  parsedEndDate: any;
  minDate: string = "";
  maxDate: string = "";
  accommID = "";
  accommodation: any | null;
  availabilityPeriods: any[] = [];
  currentUserID: string = ""
  reservationData: any | null = {
    AvailabilityPeriodId: null,
    AccommodationId: null,
    StartDate: null,
    EndDate: null,
    GuestId: null,
    GuestNum: 1,
  };

  constructor(private reservationService: ReservationService,
    private toastr: ToastrService,
    private authService: AuthService,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
    private accommodationService: AccomodationService){}
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
        this.toastr.error('Error fetching accommodation!');
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  loadAvailability(accommodationId: string){
    this.reservationService.getAvailability(accommodationId).subscribe(
      (data: any[]) => {
        this.availabilityPeriods = data;
        if (this.availabilityPeriods.length > 0) {

          const earliestStartDate = this.availabilityPeriods.reduce((min, period) =>
            period.StartDate < min ? period.StartDate : min, this.availabilityPeriods[0].StartDate);

          const latestEndDate = this.availabilityPeriods.reduce((max, period) =>
            period.EndDate > max ? period.EndDate : max, this.availabilityPeriods[0].EndDate);

          this.minDate = new Date(earliestStartDate).toISOString().split('T')[0];
          this.maxDate = new Date(latestEndDate).toISOString().split('T')[0];
        }
      },
      (error: any) => {
        this.toastr.error('Error fetching availability:');
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  reservePeriod(periodid: any, start: any, end: any): void {
    console.log(start)

    this.parsedStartDate = start + "T00:00:00Z"
    this.parsedEndDate = end + "T00:00:00Z"

    this.reservationData.AvailabilityPeriodId = periodid;
    this.reservationData.AccommodationId = this.accommID;
    this.reservationData.StartDate = this.parsedStartDate;
    this.reservationData.EndDate = this.parsedEndDate;
    this.reservationData.GuestId = this.currentUserID;

    this.reservationService.postReservation(this.reservationData)
      .subscribe(
        (response) => {
          this.toastr.success('Reservation successful');
          console.log('Reservation successful:', response);
        },
        (error) => {
          this.toastr.error('Reservation failed');
          console.error('Reservation failed:', error);

        }
      );
  }

}
