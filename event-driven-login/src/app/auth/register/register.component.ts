import { Component, computed, inject, signal } from '@angular/core';
import { RegisterService } from './register.service';
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, Validators } from '@angular/forms';
import { timer } from 'rxjs';
import { CommonModule } from '@angular/common';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIcon } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinner } from '@angular/material/progress-spinner';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {zxcvbn} from "zxcvbn-typescript";
import {
  trigger,
  transition,
  style,
  animate,
} from '@angular/animations';
import { RouterLink } from '@angular/router';


@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  imports: [CommonModule, ReactiveFormsModule, MatFormFieldModule,
    MatInputModule, MatIcon, MatCardModule,
    MatButtonModule, MatProgressSpinner, RouterLink, MatProgressBarModule],
  providers: [RegisterService],
  animations: [
    trigger('fadeAnimation', [
      transition(':enter', [
        style({ opacity: 0 }),
        animate('200ms ease-in', style({ opacity: 1 }))
      ]),
      transition(':leave', [
        animate('200ms ease-out', style({ opacity: 0 }))
      ])
    ])
  ]
})
export class RegisterComponent {
  private _registerService = inject(RegisterService);
  private _formBuilder = inject(FormBuilder);
  form = this._formBuilder.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
    password: ['', [Validators.required, this.scoreValidator]]
  });
  registerLoading = signal(false);
  registerSuccess = signal(false);
  registerError = signal(false);
  password = signal('');
  passwordSecurity = computed(() => {
    const pass = this.password();
    return (zxcvbn(pass).score / 4) * 100;
  });
  validUsername = signal(false);

  constructor() {
    this.form.valueChanges.subscribe(val => {
      this.password.set(val.password ?? '');
      this.validUsername.set(this.form.controls.username.valid);
    })
  }

  onRegister(): void {
    if (this.form.invalid) {
      return;
    }
    if (!this.form.value.username || !this.form.value.password) {
      return;
    }
    this.registerLoading.set(true);
    this.registerSuccess.set(false);
    this.registerError.set(false);
    this._registerService.register(this.form.value.username, this.form.value.password).subscribe({
      next: () => {
        this.registerSuccess.set(true);
        console.log('Registration successful');
        // TODO navigate to login
      },  error: (error) => {
        console.error('Registration failed', error);
        this.registerError.set(true);
        timer(1000).subscribe(() => {
          this.registerError.set(false);
        });
      }
    }).add(() => this.registerLoading.set(false));

  }

  scoreValidator(control: AbstractControl): ValidationErrors | null {
    if(zxcvbn(control.value).score < 2) {
      return {passwordStrength:true}
    }
    return null;
  }
}
