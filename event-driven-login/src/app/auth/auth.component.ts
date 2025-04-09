import { Component, effect, inject, signal } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { ChildrenOutletContexts, NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { fadeInAnimation } from '../route-animations';
import { animate, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  imports: [CommonModule, MatCardModule, RouterOutlet],
  animations: [fadeInAnimation,
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
export class AuthComponent {
  contexts = inject(ChildrenOutletContexts)
  pageTitle = signal('');
  visibleTitle = signal(true);

  public constructor(router: Router) {
    router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        const parts = event.url.split('/');
        const page = parts[parts.length - 1] as 'login' | 'register';
        if (page === 'login') {
          this.pageTitle.set('Login');
        } else {
          this.pageTitle.set('Register');
        }
      }
    });

    effect(() => {
      this.pageTitle();
      this.visibleTitle.set(false);
      setTimeout(() => this.visibleTitle.set(true), 300);
    })
  }

  getRouteAnimationData() {
    return this.contexts.getContext('primary')?.route?.snapshot?.data?.['animation'];
  }
}
