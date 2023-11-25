// login.component.ts
import { NgForm } from '@angular/forms';
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
  token: string|undefined;

  constructor(private authService: AuthService, private router: Router) {
    this.token = undefined;
  }

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
