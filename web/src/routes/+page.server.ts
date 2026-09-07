import { loadQuizPayload } from '$lib/server/data';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = () => loadQuizPayload();
