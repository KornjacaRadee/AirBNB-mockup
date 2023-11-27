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

    this.accomodationService
      .searchAccomodations(this.searchTerm, this.minGuests, this.maxGuests)
      .subscribe(
        (accommodations) => {
          if (accommodations && accommodations.length > 0) {
            this.accommodations = accommodations;
            this.filteredAccommodations = this.filterAccommodations();
            this.searchSuccess = true;
          } else {
            this.accommodations = [];
            this.filteredAccommodations = [];
            this.searchSuccess = false;
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
    // Implementirajte logiku za filtriranje smeštaja prema trenutnim parametrima pretrage
    // Na primer, možete koristiti metodu filter ili neki drugi pristup
    return this.accommodations.filter((accommodation) =>
      this.accommodationMatchesSearch(accommodation)
    );
  }

  accommodationMatchesSearch(accommodation: any): boolean {
    // Implementirajte logiku za proveru da li smeštaj odgovara trenutnoj pretrazi
    // Na primer, možete proveriti da li smeštaj sadrži searchTerm u svom imenu ili lokaciji
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
