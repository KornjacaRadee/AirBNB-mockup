import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
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
    return ""
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
