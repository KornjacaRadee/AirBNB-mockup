// login.component.ts
import { NgForm } from '@angular/forms';
import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent implements OnInit {
  user: any = {};
  token: string|undefined;
  errorMessage: string | undefined;

  constructor(private toastr: ToastrService,private authService: AuthService, private router: Router) {
    this.token = undefined;
  }

  ngOnInit(): void {
    if(this.authService.isAuthenticated()){
      this.router.navigate(['/home']);
    }
  }
  loginUser() {
    this.authService.login(this.user).subscribe(
      (response) => {
        console.log('Login successful', response);
        this.router.navigate(['/profile']);
      },
      (response) => {
        console.error('Login failed', response.error);
        if (response.error && response.error.includes('502')) {
          this.errorMessage = 'Service not available';
          this.toastr.error('Service not available');
        } else {
          this.errorMessage = response.error;
        }
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
