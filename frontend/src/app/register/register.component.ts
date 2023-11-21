// register.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
})
export class RegisterComponent {
  user: any = {};

  constructor(private authService: AuthService, private router: Router) {} // Dodajte Router ovde

  registerUser() {
    this.authService.register(this.user).subscribe(
      (response) => {
        console.log('Registration successful', response);
        this.router.navigate(['/']); // Prilagodite putanju prema vašoj početnoj stranici

        // Dodaj dodatne akcije po uspešnoj registraciji
      },
      (error) => {
        console.error('Registration failed', error);
        // Dodaj akcije koje ćeš preduzeti u slučaju greške
      }
    );
  }
}
