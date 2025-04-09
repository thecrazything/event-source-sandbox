import { Component, inject } from '@angular/core';
import { ChildrenOutletContexts, RouterOutlet } from '@angular/router';
import { WebSocketService } from './service/websocket.service';
import { fadeInAnimation } from './route-animations';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  animations: [fadeInAnimation]
})
export class AppComponent {
  title = 'event-driven-login';
  websocket = inject(WebSocketService);
  contexts = inject(ChildrenOutletContexts)

  public constructor() {
    this.websocket.connect();
  }


}
