import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ApiService } from '../api.service';

type Tone = 'info' | 'success' | 'error';

interface AuthorEnrollmentDraft {
  displayName: string;
  tagline: string;
  bio: string;
  experience: string;
  genres: string;
  website: string;
  socials: string;
}

@Component({
  selector: 'app-author-enrollment',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './author-enrollment.component.html'
})
export class AuthorEnrollmentComponent implements OnInit {
  step = 1;
  saving = false;
  message: { text: string; tone: Tone } | null = null;
  draftKey = 'books-author-enrollment';
  form: AuthorEnrollmentDraft = {
    displayName: '',
    tagline: '',
    bio: '',
    experience: '',
    genres: '',
    website: '',
    socials: ''
  };

  constructor(private api: ApiService, private router: Router) {}

  ngOnInit(): void {
    this.loadDraft();
  }

  nextStep(): void {
    if (this.step === 1 && !this.form.displayName.trim()) {
      this.setMessage('We need a display name to introduce you to readers.', 'error');
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
      this.setMessage('Draft saved. You can come back anytime to finish.', 'success');
    } catch {
      this.setMessage('Unable to save draft in this browser.', 'error');
    }
  }

  clearDraft(): void {
    localStorage.removeItem(this.draftKey);
  }

  submit(): void {
    if (!this.form.displayName.trim()) {
      this.setMessage('Add your display name before submitting.', 'error');
      this.step = 1;
      return;
    }
    this.saving = true;
    const sections = [
      this.form.tagline,
      this.form.bio,
      this.form.experience ? `Experience: ${this.form.experience}` : '',
      this.form.genres ? `Genres: ${this.form.genres}` : '',
      this.form.website ? `Website: ${this.form.website}` : '',
      this.form.socials ? `Socials: ${this.form.socials}` : ''
    ].filter(Boolean);
    const combinedBio = sections.join('\n\n');
    this.api.applyAuthor(this.form.displayName.trim(), combinedBio || undefined).subscribe({
      next: () => {
        this.saving = false;
        this.clearDraft();
        this.setMessage('Thanks! We received your author enrollment. We will notify you once it is live.', 'success');
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
      displayName: '',
      tagline: '',
      bio: '',
      experience: '',
      genres: '',
      website: '',
      socials: ''
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
      // ignore malformed drafts
    }
  }

  private setMessage(text: string, tone: Tone): void {
    this.message = { text, tone };
  }
}
