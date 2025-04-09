import {
  trigger,
  transition,
  style,
  animate,
  query,
  group,
  animateChild
} from '@angular/animations';

export const fadeInAnimation =
  trigger('routeAnimations', [
    transition('LoginPage <=> RegisterPage', [
      style({ position: 'relative' }),
      query(':enter, :leave', [
        style({
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          opacity: 1
        })
      ]),
      query(':enter', [
        style({ opacity: 0 })
      ], { optional: true }),
      query(':leave', animateChild(), { optional: true }),
      group([
        query(':leave', [
          animate('200ms ease-out', style({ opacity: 0 }))
        ], { optional: true }),
        query(':enter', [
          style({ opacity: 0 }),
          animate('200ms 200ms ease-out', style({ opacity: 1 })) // Add delay for fade-in
        ], { optional: true }),
      ]),
    ]),
    transition('* <=> *', [
      style({ position: 'relative' }),
      query(':enter, :leave', [
        style({
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          opacity: 1
        })
      ], { optional: true }),
      query(':enter', [
        style({ opacity: 0 })
      ], { optional: true }),
      query(':leave', [
        animate('100ms ease-out', style({ opacity: 0 }))
      ], { optional: true }),
      query(':enter', [
        style({ opacity: 0 }),
        animate('200ms 100ms ease-out', style({ opacity: 1 })) // Add delay for fade-in
      ], { optional: true }),
      query('@*', animateChild(), { optional: true })
    ])
  ]);
