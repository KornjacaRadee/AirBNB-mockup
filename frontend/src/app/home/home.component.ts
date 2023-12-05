import { Component } from '@angular/core';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { User } from '../models/User';
import { Accommodation } from '../models/Accommodation';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent {
  title = 'frontend';
  opened = false;
  searchTerm: string = '';
  minGuests: number = 0;
  maxGuests: number = 0;
  enterPressed = false;
  searchButtonPressed = false;
  accommodations: any[] = [];
  allAccoms: Accommodation[] = [];
  filteredAccommodations: any[] = [];
  searchSuccess = false;

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
    this.enterPressed = true;

    console.log('Search Term:', this.searchTerm);
    console.log('Min Guests:', this.minGuests);
    console.log('Max Guests:', this.maxGuests);

    this.accommodationService
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
        (error: any) => {
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
    const filtered = this.accommodations.filter((accommodation) =>
      this.accommodationMatchesSearch(accommodation)
    );
    console.log('Filtered Accommodations:', filtered); // Dodajte ovu liniju
    return filtered;
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
