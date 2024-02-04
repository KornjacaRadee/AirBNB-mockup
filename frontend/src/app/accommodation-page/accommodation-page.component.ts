import { Component, OnInit,  AfterViewInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { ReservationService } from '../services/reservation/reservation.service';
import { DatePipe } from '@angular/common';

import { ToastrService } from 'ngx-toastr';
import { PrimeNGConfig } from 'primeng/api';
import { each } from 'jquery';

@Component({
  selector: 'app-accommodation-page',
  templateUrl: './accommodation-page.component.html',
  styleUrls: ['./accommodation-page.component.css']
})
export class AccommodationPageComponent implements OnInit,AfterViewInit  {
  selectedStartDate: string | null = "";
  selectedEndDate: string | null = "";
  parsedStartDate: any;
  parsedEndDate: any;
  minDate: string = "";
  maxDate: string = "";
  accommID = "";
  pictures: any[] = [];
  picsdata: any[] = [];
  accommodation: any | null;
  $: any;
  availabilityPeriods: any[] = [];
  currentUserID: string = ""

  reservationData: any | null = {
    AvailabilityPeriodId: null,
    AccommodationId: null,
    StartDate: null,
    EndDate: null,
    GuestId: null,
    HostId: null,
    GuestNum: 1,
  };

  constructor(private reservationService: ReservationService,
    private toastr: ToastrService,
    private authService: AuthService,
    private router: Router,
    private route: ActivatedRoute,
    private datePipe: DatePipe,
    private config: PrimeNGConfig,
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
        this.loadPictures(this.accommID)

    });
  }
  ngAfterViewInit(): void {
    $(document).ready(() => {
      $('.carousel').each(function () {
        const $carousel = $(this);
        const $items = $carousel.find('.carousel-item');
        let currentIndex = 0;

        // Show first slide
        $items.eq(currentIndex).addClass('active');

        // Next button click handler
        $carousel.find('.carousel-control.next').click(function (e) {
          e.preventDefault();
          currentIndex = (currentIndex + 1) % $items.length;
          updateCarousel();
        });

        // Previous button click handler
        $carousel.find('.carousel-control.prev').click(function (e) {
          e.preventDefault();
          currentIndex = (currentIndex - 1 + $items.length) % $items.length;
          updateCarousel();
        });

        // Update carousel function
        function updateCarousel() {
          $items.removeClass('active');
          $items.eq(currentIndex).addClass('active');
        }
      });
    });
  }


  isHost(): boolean{
    if(this.authService.getUserRole() == "host"){

      return true
    }else{
      return false
    }
  }

  loadPictures(id: string){
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
    this.reservationData.HostId = this.accommodation.owner.id
    console.log(this.reservationData.HostId)
    console.log(this.accommodation.owner.id)



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
