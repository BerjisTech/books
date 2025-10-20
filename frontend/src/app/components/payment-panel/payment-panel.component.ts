import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-payment-panel',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './payment-panel.component.html'
})
export class PaymentPanelComponent {
  // Accept either the raw { intent, next_action } or wrapped { payment: { intent, next_action } }
  _details: any;
  @Input() set details(v: any) { this._details = v; }
  get payment() {
    if (!this._details) return null;
    if (this._details.payment) return this._details.payment;
    if (this._details.intent || this._details.next_action) return this._details;
    return null;
  }
}
