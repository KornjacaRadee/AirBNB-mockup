import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';

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
    private toastr: ToastrService,
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
      last_name: ['', [Validators.required, Validators.pattern('[a-zA-Z ]*')]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      address: ['', Validators.required],
      role: ['host', Validators.required],
    });
  }

  registerUser(): void {
    if (this.registrationForm.valid) {
      this.authService.register(this.registrationForm.value).subscribe(
        (response) => {
          console.log('Registration successful', response);
          this.toastr.success('Registration successful, you can log in now.');
          this.router.navigate(['/login']);
          // Add additional actions on successful registration
        },
        (response) => {
          console.error('Registration failed', response.error);
          if (response.error && response.error.includes('502')) {
            this.errorMessage = 'Service not available';
            this.toastr.error('Service not available');
          } else {
            this.errorMessage = response.error;
          }

          // Add actions on registration failure
        }
      );
    }
  }
}
