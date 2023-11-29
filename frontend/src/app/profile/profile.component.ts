import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { User } from '../model/user';
import { HttpErrorResponse } from '@angular/common/http';


@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit{
  user: User = {} as User ;
  constructor(private authService: AuthService) { }

  ngOnInit(): void {
    this.getUserInfo()

  }

  logout(): void {
    this.authService.logout();
  }

  getUserInfo() {

    this.authService.whoami().subscribe(
      (response: User) => {
        this.user = response;
      },
      (error: HttpErrorResponse) => {
        alert(error.message);
      }
    );
  }


}
