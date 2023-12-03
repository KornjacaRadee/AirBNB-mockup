import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
})
export class RegisterComponent {
  user: any = {};
  token: string | undefined;
  errorMessage: string | undefined;
  registrationForm!: FormGroup;

  constructor(
    private authService: AuthService,
    private router: Router,
    private formBuilder: FormBuilder
  ) {
    this.token = undefined;
  }

  ngOnInit(): void {
    this.buildForm();
    if(this.authService.isAuthenticated()){
      this.router.navigate(['/home']);
    }
  }

  buildForm(): void {
    this.registrationForm = this.formBuilder.group({
      name: ['', [Validators.required, Validators.pattern('[a-zA-Z ]*')]],
      lastName: ['', [Validators.required, Validators.pattern('[a-zA-Z ]*')]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      address: ['', Validators.required],
      role: ['host', Validators.required],
      // Add other form controls and validations as needed
      // ...
    });
  }

  registerUser(): void {
    if (this.registrationForm.valid) {
      this.authService.register(this.registrationForm.value).subscribe(
        (response) => {
          console.log('Registration successful', response);
          this.router.navigate(['/login']);
          // Add additional actions on successful registration
        },
        (response) => {
          console.error('Registration failed', response.error);
          this.errorMessage = response.error;
          // Add actions on registration failure
        }
      );
    }
  }
}
