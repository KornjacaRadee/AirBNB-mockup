import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  title = 'frontend';
  opened = false;
  searchTerm: string = '';
  enterPressed = false;
  searchButtonPressed = false;

  onSearch() {
    this.enterPressed = true;
    console.log('Searching for:', this.searchTerm);
  }

  onSearchButtonClicked() {
    this.searchButtonPressed = true;
  }
}
