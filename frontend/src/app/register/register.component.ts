// register.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
})
export class RegisterComponent {
  user: any = {};
  token: string|undefined;

  constructor(private authService: AuthService, private router: Router) {
    this.token = undefined;
  } // Dodajte Router ovde

  registerUser() {
    this.authService.register(this.user).subscribe(
      (response) => {
        console.log('Registration successful', response);
        this.router.navigate(['/login']); // Prilagodite putanju prema vašoj početnoj stranici

        // Dodaj dodatne akcije po uspešnoj registraciji
      },
      (error) => {
        console.error('Registration failed', error);
        // Dodaj akcije koje ćeš preduzeti u slučaju greške
      }
    );
  }


  // public send(form: NgForm): void {
  //   if (form.invalid) {
  //     for (const control of Object.keys(form.controls)) {
  //       form.controls[control].markAsTouched();
  //     }
  //     return;                                              //VEZANO ZA CAPTCHU, OSTAVITI ZA SVAKI SLUCAJ
  //   }

  //   console.debug(`Token [${this.token}] generated`);
  // }



}
