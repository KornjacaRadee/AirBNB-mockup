import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit{
  user: any;
  id: string = "";
  constructor(private authService: AuthService,private router: Router) { }

  ngOnInit(): void {
    if(!this.authService.isAuthenticated()){
      this.router.navigate(['/login']);
    }
    this.loadUserDetails();
    console.log(this.user);
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
        console.log('', response);
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
        console.log(response)
        // Map the response to the 'user' property
        this.user = response;
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }
}
