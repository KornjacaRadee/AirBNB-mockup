import { Component, OnInit } from '@angular/core';
import { AccomodationService } from '../services/accomodation/accomodation.service';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  accommodations: any[] = [];

  constructor(private accommodationService: AccomodationService) {}

  ngOnInit(): void {
    this.getAccommodations();
  }

  getAccommodations(): void {
    this.accommodationService.getAccomodations().subscribe(
      (data: any[]) => {
        this.accommodations = data;
      },
      (error) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }
}