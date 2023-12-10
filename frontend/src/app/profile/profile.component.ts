import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { AccomodationService } from '../services/accomodation/accomodation.service';
import { CommonModule } from '@angular/common';
@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit{
  user: any;
  id: string = "";
  accomms: any[] = [];
  showAccommodations = false;



  constructor(private authService: AuthService,private accommodationService: AccomodationService,private router: Router) { }

  ngOnInit(): void {

    if(!this.authService.isAuthenticated()){
      this.router.navigate(['/login']);
    }
    this.accomms = [];
    this.getUserAccommodations()
    this.loadUserDetails();
    console.log(this.accomms);
  }
  getUserAccommodations(): void {
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

    this.accommodationService.getUserAccommodations(headers).subscribe(
      (data: any[]) => {
        console.log(data)
        this.accomms = data;
      },
      (error: any) => {
        console.error('Error fetching accommodations:', error);
      }
    );
  }

  logout(): void {
    this.authService.logout();
  }

  editProfile(){
    return ""
  }

  deleteProfile(){
    const token = this.authService.getAuthToken();

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);

    this.authService.deleteUser(headers).subscribe(
      (response) => {
        this.authService.logout();
      },
      (error) => {
        // Handle error, e.g., show an error message
        console.error('Failed to delete user', error);
      }
    );


}



  loadUserDetails() {
    this.id = this.authService.getUserId();
    this.authService.getUserById(this.id).subscribe(
      (response) => {
        // Map the response to the 'user' property
        this.user = response;
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }
  toggleAccommodations() {
    this.showAccommodations = !this.showAccommodations;
  }
}
