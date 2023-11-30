import { Component } from '@angular/core';
import { AccomodationService } from './services/accomodation/accomodation.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  title = 'frontend';
  opened = false;
  searchTerm: string = '';
  minGuests: number = 0;
  maxGuests: number = 0;
  enterPressed = false;
  searchButtonPressed = false;
  accommodations: any[] = [];
  filteredAccommodations: any[] = [];
  searchSuccess = false;

  constructor(private accomodationService: AccomodationService) {}

  onSearch() {
    this.enterPressed = true;

    console.log('Search Term:', this.searchTerm);
    console.log('Min Guests:', this.minGuests);
    console.log('Max Guests:', this.maxGuests);

    this.accomodationService
      .searchAccomodations(this.searchTerm, this.minGuests, this.maxGuests)
      .subscribe(
        (response) => {
          console.log('Server Response:', response);

          if (response && response.length > 0) {
            this.accommodations = response;
            this.filteredAccommodations = this.filterAccommodations();
            this.searchSuccess = true;
            console.log(
              'Filtered Accommodations:',
              this.filteredAccommodations
            );
          } else {
            this.accommodations = [];
            this.filteredAccommodations = [];
            this.searchSuccess = false;
            console.log('No accommodations found');
          }
        },
        (error) => {
          console.error('Error searching accommodations:', error);
          this.accommodations = [];
          this.filteredAccommodations = [];
          this.searchSuccess = false;
        }
      );
  }

  onSearchButtonClicked() {
    this.searchButtonPressed = true;
    this.onSearch();
  }

  resetSearch() {
    this.accommodations = [];
    this.filteredAccommodations = [];
    this.searchSuccess = false;
    this.searchTerm = '';
    this.minGuests = 0;
    this.maxGuests = 0;
  }

  filterAccommodations(): any[] {
    return this.accommodations.filter((accommodation) =>
      this.accommodationMatchesSearch(accommodation)
    );
  }

  accommodationMatchesSearch(accommodation: any): boolean {
    return (
      (accommodation.name
        .toLowerCase()
        .includes(this.searchTerm.toLowerCase()) ||
        accommodation.location
          .toLowerCase()
          .includes(this.searchTerm.toLowerCase())) &&
      accommodation.minGuestNum >= this.minGuests &&
      accommodation.maxGuestNum <= this.maxGuests
    );
  }
}
