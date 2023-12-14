import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { ReservationService } from '../services/reservation/reservation.service';

@Component({
  selector: 'app-create-period',
  templateUrl: './create-period.component.html',
  styleUrls: ['./create-period.component.css']
})
export class CreatePeriodComponent implements OnInit {
  startDate: string = "";
  endDate: string = "";
  accommID: string = "";
  price: number = 0;
  isPricePerGuest: boolean = false;

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.accommID = params['id'];

    });
  }
  constructor(private http: HttpClient,private router: Router,private route: ActivatedRoute, private reservationService: ReservationService) {}

  createAvailabilityPeriod() {
    const startDatea = new Date(this.startDate);
    const formattedStartDate = startDatea.toISOString();
    const endDatea = new Date(this.endDate);
    const formattedEndtDate = endDatea.toISOString();
    const availabilityPeriod = {
      AccommodationId: this.accommID,
      StartDate: formattedStartDate,
      EndDate: formattedEndtDate,
      Price: this.price,
      IsPricePerGuest: this.isPricePerGuest
    };

    this.reservationService.postAvailability(availabilityPeriod).subscribe(
      (response) => {
        console.log('Availability period created successfully:', response);
        // Handle success, e.g., show a success message
      },
      (error) => {
        console.error('Error creating availability period:', error);
        // Handle error, e.g., show an error message
      }
    );
  }
}
