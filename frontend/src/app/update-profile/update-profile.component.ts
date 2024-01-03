// update-profile.component.ts
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ProfilesService } from '../services/profile/profiles.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-update-profile',
  templateUrl: './update-profile.component.html',
  styleUrls: ['./update-profile.component.css'],
})
export class UpdateProfileComponent implements OnInit {
  profile: any = {};
  email: string = '';
  updatedProfile: any = {
    address: '',
    email: '',
    username: '',
  };

  constructor(
    private toastr: ToastrService,
    private route: ActivatedRoute,
    private router: Router,
    private profileService: ProfilesService
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      this.email = params['email'];
      this.loadProfileDetails();
    });
  }

  loadProfileDetails() {
    this.profileService.getProfileByEmail(this.email).subscribe(
      (response) => {
        // Map the response to the 'user' property
        this.profile = response;
        console.log(this.profile);
      },
      (error) => {
        console.error('Error fetching user details', error);
      }
    );
  }

  saveProfile(): void {
    this.profileService
      .updateProfileByEmail(this.profile.id, this.profile)
      .subscribe(
        () => {
          this.router.navigate(['/profile']);
          this.toastr.success('Profile saved successfully');
        },
        (error: any) => {
          this.toastr.error('Error editing profile! Try again. :)');
          console.error('Error:', error);
        }
      );
  }
}
