import { Component, Input, Output, EventEmitter } from '@angular/core';
import { Store } from '../../services/models/stores-response.model';
import { determineLevel } from '../../services/utils';

@Component({
  selector: 'app-store-card',
  templateUrl: './store-card.component.html',
  styleUrls: ['./store-card.component.scss']
})
export class StoreCardComponent {
  @Input() data: Store;
  @Output() moveToMedical: EventEmitter<Store> = new EventEmitter();

  moveMap(data: Store) {
    this.moveToMedical.emit(data);
  }

  isExpirationNotExpired(date: string | Date) {
    return new Date(date).valueOf() > 0;
  }
}
