// login.component.ts

import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent {
  user: any = {};

  constructor(private authService: AuthService, private router: Router) {}

  loginUser() {
    this.authService.login(this.user).subscribe(
      (response) => {
        console.log('Login successful', response);
        this.router.navigate(['/home']);
      },
      (error) => {
        console.error('Login failed', error);
        // Dodaj akcije koje ćeš preduzeti u slučaju greške
      }
    );
  }
}
