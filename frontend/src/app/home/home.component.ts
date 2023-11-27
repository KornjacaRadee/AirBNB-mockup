import { Component, OnInit } from '@angular/core';
import { AccomodationService } from '../services/accomodation/accomodation.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  accommodations: any[] = [];
  searchTerm: string = '';

  constructor(private accommodationService: AccomodationService) {}

  ngOnInit(): void {
    this.getAccommodations();
  }

  getAccommodations(): void {
    this.accommodationService.getAccomodations().subscribe(
      (data: any[]) => {
        this.accommodations = data;
      },
      (error: any) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  onSearch() {
    this.accommodations = this.accommodations.filter((accommodation) =>
      this.accommodationMatchesSearch(accommodation)
    );
  }

  accommodationMatchesSearch(accommodation: any): boolean {
    return (
      accommodation.name
        .toLowerCase()
        .includes(this.searchTerm.toLowerCase()) ||
      accommodation.location
        .toLowerCase()
        .includes(this.searchTerm.toLowerCase())
    );
  }
}
