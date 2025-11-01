import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';

type Tone = 'info' | 'success' | 'error';

interface PublisherEnrollmentDraft {
  name: string;
  imprint: string;
  description: string;
  focus: string;
  distribution: string;
  website: string;
  contact: string;
}

@Component({
  selector: 'app-publisher-enrollment',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './publisher-enrollment.component.html'
})
export class PublisherEnrollmentComponent implements OnInit {
  step = 1;
  saving = false;
  message: { text: string; tone: Tone } | null = null;
  draftKey = 'books-publisher-enrollment';
  form: PublisherEnrollmentDraft = {
    name: '',
    imprint: '',
    description: '',
    focus: '',
    distribution: '',
    website: '',
    contact: ''
  };

  constructor(private api: ApiService, private router: Router) {}

  ngOnInit(): void {
    this.loadDraft();
  }

  nextStep(): void {
    if (this.step === 1 && !this.form.name.trim()) {
      this.setMessage('Let us know the name of your publishing house or imprint.', 'error');
      return;
    }
    this.step = Math.min(3, this.step + 1);
  }

  prevStep(): void {
    this.step = Math.max(1, this.step - 1);
  }

  saveDraft(): void {
    try {
      localStorage.setItem(this.draftKey, JSON.stringify(this.form));
      this.setMessage('Draft saved locally.', 'success');
    } catch {
      this.setMessage('Unable to save draft in this browser.', 'error');
    }
  }

  clearDraft(): void {
    localStorage.removeItem(this.draftKey);
  }

  submit(): void {
    if (!this.form.name.trim()) {
      this.setMessage('Add your publisher name before submitting.', 'error');
      this.step = 1;
      return;
    }
    this.saving = true;
    const sections = [
      this.form.imprint ? `Imprint: ${this.form.imprint}` : '',
      this.form.description,
      this.form.focus ? `Focus areas: ${this.form.focus}` : '',
      this.form.distribution ? `Distribution: ${this.form.distribution}` : '',
      this.form.website ? `Website: ${this.form.website}` : '',
      this.form.contact ? `Contact: ${this.form.contact}` : ''
    ].filter(Boolean);
    const combined = sections.join('\n\n');
    this.api.applyPublisher(this.form.name.trim(), combined || undefined).subscribe({
      next: () => {
        this.saving = false;
        this.clearDraft();
        this.setMessage('Thanks! Your publisher enrollment is in review.', 'success');
        this.step = 3;
      },
      error: () => {
        this.saving = false;
        this.setMessage('We could not submit the enrollment right now. Please try again shortly.', 'error');
      }
    });
  }

  startOver(): void {
    this.clearDraft();
    this.form = {
      name: '',
      imprint: '',
      description: '',
      focus: '',
      distribution: '',
      website: '',
      contact: ''
    };
    this.step = 1;
    this.message = null;
  }

  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }

  private loadDraft(): void {
    try {
      const stored = localStorage.getItem(this.draftKey);
      if (stored) {
        const parsed = JSON.parse(stored);
        this.form = { ...this.form, ...parsed };
        this.setMessage('Draft loaded from your last visit.', 'info');
      }
    } catch {
      // ignore
    }
  }

  private setMessage(text: string, tone: Tone): void {
    this.message = { text, tone };
  }
}
