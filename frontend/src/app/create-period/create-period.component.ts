import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { ReservationService } from '../services/reservation/reservation.service';
import { ToastrService } from 'ngx-toastr';

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
  constructor( private toastr: ToastrService,private http: HttpClient,private router: Router,private route: ActivatedRoute, private reservationService: ReservationService) {}

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
        this.toastr.success('Period created successfully');
        console.log('Availability period created successfully:', response);
      },
      (error) => {
        this.toastr.error('Error creating period');
        console.error('Error creating availability period:', error);
      }
    );
  }
}
